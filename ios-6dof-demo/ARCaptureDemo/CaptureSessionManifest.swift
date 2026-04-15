import Foundation

struct CaptureSessionManifest: Codable {
    let sessionId: String
    let startedAtISO8601: String
    let hasDepthSensor: Bool
    let frameIntervalSeconds: Double
    let records: [CaptureRecord]
}
