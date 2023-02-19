import {Button, Space, Table, Input, message} from "antd";
import {CloudUploadOutlined, DeleteOutlined} from "@ant-design/icons";
import React, {useCallback, useEffect, useState} from "react";
import {ColumnsType} from "antd/es/table";
import {DownloadFile, LocalIPAddr} from "../../../wailsjs/go/main/App"
import {LogInfo} from "../../../wailsjs/runtime";
import {isEqualLAN} from "../../utils";
import {decrypt, encrypt} from "../../utils/crypto";
const {Search} = Input
interface DataType {
    key: React.Key;
    file: string;
    path: string;
    size: number;
}

const columns: ColumnsType<DataType> = [
    {
        title: '文件名',
        dataIndex: 'file',
    },
    {
        title: '文件路径',
        dataIndex: 'path',
    },
    {
        title: '大小/kb',
        dataIndex: 'size',
    },
];

export default function Download() {
    const [files, setFiles] = useState<Array<DataType>>([]);
    const [selectedFiles, setSelectedFiled] = useState<Array<DataType>>([]);
    const [localIp, setLoaclIp] = useState<string[]>([])
    const onDownloadFiles = useCallback(async () => {
        let [ip, port] = await LocalIPAddr()
        // let res = await DownloadFile('/Users/chenlong/Documents/述职报告--陈龙.pages', "/Users/chenlong/Documents/述职报告--陈龙1.pages")
        console.log(ip, port)
        // console.log(res)

    }, []);
    useEffect(() => {
        LocalIPAddr().then(([ip, port]) => {
            setLoaclIp( [ip, port])
        })
        return
    }, [])
    const onDeleteFiles = useCallback(() => {

    }, [])
    const onSearch = (value: string) => {
        if(!value) {
            return message.error("请传入hash")
        }
        console.log(encrypt(localIp.join(',')))
        console.log(decrypt(encrypt(localIp.join(','))))
        // const [ip, port] =
        if(!isEqualLAN(localIp[0], "192.168.11.1")) {
            return message.error("发送端IP与本机IP不属于同一局域网")
        }
        console.log(localIp)

    }
  return (
      <div className='download'>
              <Search placeholder="请输入传输密钥" onSearch={onSearch} style={{
                  width: 'calc(100% - 48px)',
                  margin: '16px 0'
              }}></Search>
          <Space>
              <Button icon={<CloudUploadOutlined />} size={'middle'} onClick={onDownloadFiles}>
                  下载
              </Button>
              <Button icon={<DeleteOutlined />} size={'middle'} onClick={onDeleteFiles}>
                  移除
              </Button>
          </Space>
          <Table
              rowSelection={{
                  type: 'checkbox',
                  onChange: (selectedRowKeys: React.Key[], selectedRows: DataType[]) => {
                      setSelectedFiled(selectedRows);
                  },
                  getCheckboxProps: (record: DataType) => ({
                      disabled: false,
                      ...record,
                  }),
              }}
              columns={columns}
              dataSource={files}
              pagination={false}
              scroll={{ y: 500 }}
              style={{
                  marginTop: 15,
              }}
          />
      </div>
  );
}
