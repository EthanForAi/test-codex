import SwiftUI

struct ContentView: View {
    @StateObject private var manager = ARCaptureManager()

    var body: some View {
        ZStack(alignment: .bottom) {
            ARSessionView(session: manager.session)
                .ignoresSafeArea()

            VStack(spacing: 12) {
                Text(manager.statusText)
                    .font(.footnote)
                    .foregroundStyle(.white)
                    .padding(.horizontal, 10)
                    .padding(.vertical, 6)
                    .background(.black.opacity(0.55), in: Capsule())

                HStack(spacing: 12) {
                    Button("Start Capture") {
                        manager.startCapture()
                    }
                    .buttonStyle(.borderedProminent)
                    .tint(.green)
                    .disabled(manager.isCapturing)

                    Button("Stop Capture") {
                        manager.stopCapture()
                    }
                    .buttonStyle(.borderedProminent)
                    .tint(.red)
                    .disabled(!manager.isCapturing)
                }
            }
            .padding(.bottom, 24)
        }
        .onAppear {
            manager.startSessionIfNeeded()
        }
    }
}

#Preview {
    ContentView()
}
