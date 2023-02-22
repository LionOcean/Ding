import { useNavigate } from 'react-router-dom';
import { Popconfirm } from 'antd';
import { WindowMinimise, Quit } from '../../../wailsjs/runtime';
import { CloseCircleOutlined, HomeOutlined, MinusCircleOutlined } from '@ant-design/icons';
import './index.less';

export default function () {
  const navigate = useNavigate();

  return (
    <div className='header'>
      <div className='header__backBtn' onClick={() => navigate('/home')} role='button' tabIndex={-1}>
        <HomeOutlined style={{ fontSize: '20px' }} />
      </div>
      <span
        className='header__title'
        style={{
          ['--wails-draggable' as string]: 'drag',
        }}
      >
        Ding
      </span>
      <div className='header__operation'>
        <MinusCircleOutlined style={{ fontSize: '18px', marginRight: '10px' }} onClick={() => WindowMinimise()} role='button' tabIndex={-1} />
        <Popconfirm title='退出' description='你确定要现在退出Ding嘛?' placement='bottomRight' okText='确定' cancelText='不' onConfirm={Quit}>
          <CloseCircleOutlined style={{ fontSize: '18px' }} role='button' tabIndex={-1} />
        </Popconfirm>
      </div>
    </div>
  );
}
