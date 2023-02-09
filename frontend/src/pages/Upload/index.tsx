import {FileOutlined, CloudUploadOutlined } from "@ant-design/icons"
import {useState} from "react";
import { Card} from "antd"
import "./index.less"
export default function Upload() {
    interface File {
        name: string,
        path: string,
        [key:string]: any
    }
    const [fileList, setFileList] = useState<File[]>([{
        name: "文件",
        path: "/User/Local/bin/写真.text"
    }])
    const uploadFile = (file:File) => {

    }
    return (
    <div className={"file-list"}>
        { fileList?.map(file => (
            <Card className={'file-item'} key={file.name}>
                <FileOutlined className={'file-icon'}/>
                <div className={'file-desc'}>
                    <div className={'file-name'}>{file.name}</div>
                    <div className={'file-path'}>{ file.path }</div>
                </div>
                <CloudUploadOutlined className={'file-icon upload-icon'} onClick={() => uploadFile(file)} />
            </Card>
        ))
        }
    </div>
  );
}
