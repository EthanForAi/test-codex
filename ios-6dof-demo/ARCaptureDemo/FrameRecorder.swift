import Foundation
import ARKit
import UIKit

final class FrameRecorder {
    private(set) var isRecording = false
    private(set) var sessionDirectoryURL: URL?

    private let frameInterval: TimeInterval
    private var lastSavedTimestamp: TimeInterval = 0
    private var records: [CaptureRecord] = []
    private var sessionId = ""
    private var hasDepthSensor = false
    private var sessionStartDate: Date?

    init(frameInterval: TimeInterval = 0.5) {
        self.frameInterval = frameInterval
    }

    func startNewSession(hasDepthSensor: Bool) throws {
        let fm = FileManager.default
        let capturesRoot = try documentsDirectory()
            .appendingPathComponent("captures", isDirectory: true)
        try fm.createDirectory(at: capturesRoot, withIntermediateDirectories: true)

        sessionId = "session_\(Int(Date().timeIntervalSince1970))"
        let sessionURL = capturesRoot.appendingPathComponent(sessionId, isDirectory: true)
        try fm.createDirectory(at: sessionURL, withIntermediateDirectories: true)

        self.hasDepthSensor = hasDepthSensor
        sessionDirectoryURL = sessionURL
        sessionStartDate = Date()
        lastSavedTimestamp = 0
        records = []
        isRecording = true
    }

    func appendFrame(_ frame: ARFrame) {
        guard isRecording, let sessionDirectoryURL else { return }

        let timestamp = frame.timestamp
        if (timestamp - lastSavedTimestamp) < frameInterval {
            return
        }

        guard let jpegData = makeJPEG(from: frame.capturedImage) else {
            return
        }

        let index = records.count
        let imageFile = String(format: "image_%05d.jpg", index)
        let imageURL = sessionDirectoryURL.appendingPathComponent(imageFile)

        do {
            try jpegData.write(to: imageURL)
        } catch {
            print("[FrameRecorder] Failed writing JPEG: \(error)")
            return
        }

        var depthFileName: String?
        var depthWidth: Int?
        var depthHeight: Int?
        var depthBytesPerRow: Int?

        if let depthMap = frame.sceneDepth?.depthMap,
           let depthRaw = makeDepthRawBytes(from: depthMap) {
            let name = String(format: "depth_%05d.raw", index)
            let url = sessionDirectoryURL.appendingPathComponent(name)
            do {
                try depthRaw.write(to: url)
                depthFileName = name
                depthWidth = CVPixelBufferGetWidth(depthMap)
                depthHeight = CVPixelBufferGetHeight(depthMap)
                depthBytesPerRow = CVPixelBufferGetBytesPerRow(depthMap)
            } catch {
                print("[FrameRecorder] Failed writing depth map: \(error)")
            }
        }

        records.append(
            CaptureRecord(
                imageFile: imageFile,
                timestamp: timestamp,
                transform: flattenTransform(frame.camera.transform),
                depthFile: depthFileName,
                depthWidth: depthWidth,
                depthHeight: depthHeight,
                depthBytesPerRow: depthBytesPerRow
            )
        )

        lastSavedTimestamp = timestamp
    }

    func stopAndFlushManifest() throws -> URL? {
        guard let sessionDirectoryURL else { return nil }

        let manifest = CaptureSessionManifest(
            sessionId: sessionId,
            startedAtISO8601: ISO8601DateFormatter().string(from: sessionStartDate ?? Date()),
            hasDepthSensor: hasDepthSensor,
            frameIntervalSeconds: frameInterval,
            records: records
        )

        let jsonData = try JSONEncoder.prettyPrinted.encode(manifest)
        let jsonURL = sessionDirectoryURL.appendingPathComponent("metadata.json")
        try jsonData.write(to: jsonURL)

        let finishedURL = sessionDirectoryURL
        isRecording = false
        self.sessionDirectoryURL = nil
        sessionId = ""
        hasDepthSensor = false
        sessionStartDate = nil
        return finishedURL
    }

    private func flattenTransform(_ matrix: simd_float4x4) -> [Float] {
        [
            matrix.columns.0.x, matrix.columns.0.y, matrix.columns.0.z, matrix.columns.0.w,
            matrix.columns.1.x, matrix.columns.1.y, matrix.columns.1.z, matrix.columns.1.w,
            matrix.columns.2.x, matrix.columns.2.y, matrix.columns.2.z, matrix.columns.2.w,
            matrix.columns.3.x, matrix.columns.3.y, matrix.columns.3.z, matrix.columns.3.w
        ]
    }

    private func makeJPEG(from pixelBuffer: CVPixelBuffer) -> Data? {
        let ciImage = CIImage(cvPixelBuffer: pixelBuffer)
        let context = CIContext(options: nil)

        guard let cgImage = context.createCGImage(ciImage, from: ciImage.extent) else {
            return nil
        }

        let image = UIImage(cgImage: cgImage)
        return image.jpegData(compressionQuality: 0.9)
    }

    private func makeDepthRawBytes(from depthMap: CVPixelBuffer) -> Data? {
        CVPixelBufferLockBaseAddress(depthMap, .readOnly)
        defer { CVPixelBufferUnlockBaseAddress(depthMap, .readOnly) }

        guard let baseAddress = CVPixelBufferGetBaseAddress(depthMap) else {
            return nil
        }

        let height = CVPixelBufferGetHeight(depthMap)
        let bytesPerRow = CVPixelBufferGetBytesPerRow(depthMap)
        let dataSize = height * bytesPerRow
        return Data(bytes: baseAddress, count: dataSize)
    }

    private func documentsDirectory() throws -> URL {
        guard let url = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask).first else {
            throw NSError(domain: "FrameRecorder", code: 1, userInfo: [NSLocalizedDescriptionKey: "Documents directory not found"])
        }
        return url
    }
}

private extension JSONEncoder {
    static var prettyPrinted: JSONEncoder {
        let encoder = JSONEncoder()
        encoder.outputFormatting = [.prettyPrinted, .sortedKeys]
        return encoder
    }
}
