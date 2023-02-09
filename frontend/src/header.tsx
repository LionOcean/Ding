import { WindowMinimise, Quit } from "../wailsjs/runtime"

export default function() {
  return (
    <div style={{
      ["--wails-draggable" as string ]: "drag",
      cursor: "grab",
      position: "relative",
      padding: "5px 0"
    }}>
      <div style={{
        position: "absolute",
        left: "20px",
        top: "50%",
        transform: "translate(0, -50%)",
        cursor: "pointer"
      }} onClick={() => window.history.back()}>返回</div>
      <span>我是拖动头部</span>
      <div style={{
        position: "absolute",
        right: "20px",
        top: "50%",
        transform: "translate(0, -50%)",
        cursor: "default"
      }}>
        <span style={{ cursor: "pointer" }} onClick={() => WindowMinimise()}>最小化</span>
        <span style={{
          marginLeft: "10px",
          cursor: "pointer"
        }} onClick={() => Quit()}>退出</span>
      </div>
    </div>
  )
}