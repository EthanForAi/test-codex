import SwiftUI
import ARKit
import SceneKit

struct ARSessionView: UIViewRepresentable {
    let session: ARSession

    func makeUIView(context: Context) -> ARSCNView {
        let view = ARSCNView(frame: .zero)
        view.automaticallyUpdatesLighting = true
        view.scene = SCNScene()
        view.session = session
        return view
    }

    func updateUIView(_ uiView: ARSCNView, context: Context) {}
}
