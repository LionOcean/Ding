import { useEffect } from "react";
import { HashRouter } from 'react-router-dom';
import { WindowShow } from '../wailsjs/runtime'
import Router from './routers';

import { Layout } from 'antd';

import Header from './components/Header';
import Navigate from './components/Navigate';

import './App.less';

const { Content, Footer } = Layout;

function App() {
  useEffect(() => {
    // 等待js渲染完再加载软件窗口，有效避免启动白屏
    WindowShow()
  }, [])
  return (
    <HashRouter>
      <Layout className='app__layout'>
        <Header />
        <Content>
          <Router />
        </Content>
        <Navigate />
      </Layout>
    </HashRouter>
  );
}

// console.log('__APP_VERSION__: ', __APP_VERSION__)

export default App;
