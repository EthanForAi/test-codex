# ARCaptureDemo (SwiftUI + ARKit)

最小可运行的 6DoF 数据采集 Demo。

## 文件放置位置
在 Xcode 新建 iOS App（SwiftUI）后，把以下文件放到项目主 target（例如 `ARCaptureDemo`）中：

- `ARCaptureDemoApp.swift`
- `ContentView.swift`
- `ARSessionView.swift`
- `ARCaptureManager.swift`
- `FrameRecorder.swift`
- `CaptureRecord.swift`
- `CaptureSessionManifest.swift`

本仓库示例路径：`ios-6dof-demo/ARCaptureDemo/`。

## 必改配置（Info.plist）
增加：
- `Privacy - Camera Usage Description` (`NSCameraUsageDescription`)
  - 示例值：`需要使用相机进行 AR 数据采集。`

## 运行设备
- 需要真机（例如 iPhone 15 Pro Max）
- 需要 iOS 17+（建议）
- Simulator 无法使用 ARKit World Tracking

## 采集输出目录
程序停止采集后会把数据写入：

`Documents/captures/session_xxx/`

每个 session 包含：
- `image_00000.jpg` ... （按 0.5 秒间隔采样）
- `depth_00000.raw` ...（仅当 sceneDepth 可用时）
- `metadata.json` （图片文件名、时间戳、4x4 transform、depth 文件名）
