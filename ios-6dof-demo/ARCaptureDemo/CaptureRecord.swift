import Foundation

struct CaptureRecord: Codable {
    let imageFile: String
    let timestamp: TimeInterval
    let transform: [Float]
    let depthFile: String?
    let depthWidth: Int?
    let depthHeight: Int?
    let depthBytesPerRow: Int?
}
