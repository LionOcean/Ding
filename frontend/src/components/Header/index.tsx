import { useNavigate, useLocation } from 'react-router-dom';
import { Popconfirm } from 'antd';
import { WindowMinimise, Quit } from '../../../wailsjs/runtime';
import { CloseCircleOutlined, HomeOutlined, MinusCircleOutlined } from '@ant-design/icons';
import './index.less';

export default function () {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  // home页不展示title和返回首页按钮
  const show = pathname !== '/home';

  return (
    <div className='header'>
      <div className='header__backBtn' role='button' tabIndex={-1}>
        {show ? <HomeOutlined title='首页' onClick={() => navigate('/home')} style={{ fontSize: '20px' }} /> : ''}
      </div>
      <span
        className='header__title'
        style={{
          ['--wails-draggable' as string]: 'drag',
        }}
      >
        {show ? 'Ding' : ''}
      </span>
      <div className='header__operation'>
        <MinusCircleOutlined
          title='最小化'
          style={{ fontSize: '18px', marginRight: '10px' }}
          onClick={() => WindowMinimise()}
          role='button'
          tabIndex={-1}
        />
        <Popconfirm title='退出' description='你确定要现在退出Ding嘛?' placement='bottomRight' okText='确定' cancelText='不' onConfirm={Quit}>
          <CloseCircleOutlined title='退出' style={{ fontSize: '18px' }} role='button' tabIndex={-1} />
        </Popconfirm>
      </div>
    </div>
  );
}
