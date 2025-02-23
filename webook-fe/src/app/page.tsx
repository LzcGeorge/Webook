'use client';

import React from 'react';
import { Button, Space, Typography } from 'antd';
import { ReadOutlined, UserOutlined, EditOutlined } from '@ant-design/icons';

const { Title } = Typography;

const App = () => {
    return (
        <div style={{ 
            padding: '50px', 
            display: 'flex', 
            flexDirection: 'column',
            alignItems: 'center',
            gap: '20px'
        }}>
            <Title level={1}>小微书</Title>
            <Space direction="vertical" size="large">
                <Button 
                    type="primary" 
                    icon={<UserOutlined />} 
                    size="large"
                    href="/users/login"
                    block
                >
                    登录/注册
                </Button>
                <Button 
                    type="default" 
                    icon={<ReadOutlined />}
                    size="large"
                    href="/articles/list"
                    block
                >
                    文章列表
                </Button>
                <Button 
                    type="default" 
                    icon={<EditOutlined />}
                    size="large"
                    href="/articles/edit"
                    block
                >
                    写文章
                </Button>
            </Space>
        </div>
    );
};

export default App;