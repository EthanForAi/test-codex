import Foundation
import ARKit

final class ARCaptureManager: NSObject, ObservableObject {
    @Published var isCapturing = false
    @Published var statusText = "Initializing AR Session..."
    @Published var isSessionInterrupted = false

    let session = ARSession()

    private let recorder = FrameRecorder(frameInterval: 0.5)
    private var hasConfiguredSession = false

    override init() {
        super.init()
        session.delegate = self
    }

    func startSessionIfNeeded() {
        guard !hasConfiguredSession else { return }
        guard ARWorldTrackingConfiguration.isSupported else {
            updatePublishedStatus(text: "ARWorldTrackingConfiguration is not supported on this device.")
            return
        }

        let configuration = ARWorldTrackingConfiguration()
        configuration.worldAlignment = .gravity

        if ARWorldTrackingConfiguration.supportsFrameSemantics(.sceneDepth) {
            configuration.frameSemantics.insert(.sceneDepth)
        }

        session.run(configuration, options: [.resetTracking, .removeExistingAnchors])
        hasConfiguredSession = true

        if ARWorldTrackingConfiguration.supportsFrameSemantics(.sceneDepth) {
            updatePublishedStatus(text: "AR running (sceneDepth supported). Ready to capture.")
        } else {
            updatePublishedStatus(text: "AR running (sceneDepth unavailable). Ready to capture.")
        }
    }

    func startCapture() {
        guard !isCapturing else { return }
        startSessionIfNeeded()

        let supportsDepth = ARWorldTrackingConfiguration.supportsFrameSemantics(.sceneDepth)

        do {
            try recorder.startNewSession(hasDepthSensor: supportsDepth)
            updatePublishedStates(isCapturing: true, statusText: "Capturing...")
        } catch {
            updatePublishedStatus(text: "Failed to start capture: \(error.localizedDescription)")
        }
    }

    func stopCapture() {
        guard isCapturing else { return }

        do {
            let savedSessionURL = try recorder.stopAndFlushManifest()

            if let sessionURL = savedSessionURL {
                updatePublishedStates(isCapturing: false, statusText: "Capture saved: \(sessionURL.lastPathComponent)")
            } else {
                updatePublishedStates(isCapturing: false, statusText: "Capture stopped (no session directory).")
            }
        } catch {
            updatePublishedStates(isCapturing: false, statusText: "Stop failed: \(error.localizedDescription)")
        }
    }

    private func updatePublishedStatus(text: String) {
        DispatchQueue.main.async {
            self.statusText = text
        }
    }

    private func updatePublishedStates(isCapturing: Bool, statusText: String) {
        DispatchQueue.main.async {
            self.isCapturing = isCapturing
            self.statusText = statusText
        }
    }

    private func updateInterruptionState(isInterrupted: Bool, statusText: String) {
        DispatchQueue.main.async {
            self.isSessionInterrupted = isInterrupted
            self.statusText = statusText
        }
    }

    private func trackingLimitedReasonDescription(_ reason: ARCamera.TrackingState.Reason) -> String {
        switch reason {
        case .initializing:
            return "initializing"
        case .excessiveMotion:
            return "excessiveMotion"
        case .insufficientFeatures:
            return "insufficientFeatures"
        case .relocalizing:
            return "relocalizing"
        @unknown default:
            return "unknown"
        }
    }
}

extension ARCaptureManager: ARSessionDelegate {
    func session(_ session: ARSession, didUpdate frame: ARFrame) {
        guard isCapturing else { return }
        recorder.appendFrame(frame)
    }

    func session(_ session: ARSession, didFailWithError error: Error) {
        updatePublishedStatus(text: "ARSession failed: \(error.localizedDescription)")
    }

    func sessionWasInterrupted(_ session: ARSession) {
        updateInterruptionState(isInterrupted: true, statusText: "Session interrupted. Capture paused.")
        if isCapturing {
            stopCapture()
        }
    }

    func sessionInterruptionEnded(_ session: ARSession) {
        updateInterruptionState(isInterrupted: false, statusText: "Interruption ended. Resetting session...")
        hasConfiguredSession = false
        startSessionIfNeeded()
    }

    func session(_ session: ARSession, cameraDidChangeTrackingState camera: ARCamera) {
        switch camera.trackingState {
        case .normal:
            if isCapturing {
                updatePublishedStatus(text: "Capturing...")
            } else {
                updatePublishedStatus(text: "Tracking normal. Ready to capture.")
            }
        case .notAvailable:
            updatePublishedStatus(text: "Tracking unavailable.")
        case .limited(let reason):
            let reasonText = trackingLimitedReasonDescription(reason)
            updatePublishedStatus(text: "Tracking limited: \(reasonText)")
        }
    }
}
