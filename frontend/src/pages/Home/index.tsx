import logo from "../../assets/images/logo.png";
import { DownloadOutlined, CloudUploadOutlined } from "@ant-design/icons";
import { Card } from "antd";
import "./index.less";
import { useNavigate, Link } from "react-router-dom";
function IndexPage() {
	const navigate = useNavigate();
	return (
		<div>
			<img src={logo} id='logo' alt='logo' />
			<title>Ding</title>
			<div className={"description"}>咻的一下就过去了</div>
			<div className={"card-box"}>
				<Card className={"custom-card"} bordered={false} onClick={() => navigate("/upload")}>
					<CloudUploadOutlined />
					<div className={"text"}>上传</div>
				</Card>
				<Card className={"custom-card"} bordered={false} onClick={() => navigate("/download")}>
					<DownloadOutlined />
					<div className={"text"}>下载</div>
				</Card>
			</div>
		</div>
	);
}

export default IndexPage;
