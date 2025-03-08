'use client';

import React from 'react';
import { Button, Space, Typography, Input, Card } from 'antd';
import { ReadOutlined, UserOutlined, EditOutlined, EyeOutlined } from '@ant-design/icons';

const { Title } = Typography;

const App = () => {
    const [articleId, setArticleId] = React.useState('1');

    return (
        <div style={{ 
            minHeight: '100vh',
            background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
            padding: window.innerWidth > 768 ? '50px' : '20px', 
            display: 'flex', 
            flexDirection: 'column',
            alignItems: 'center',
            gap: '20px'
        }}>
            <Card 
                style={{ 
                    width: '100%',
                    maxWidth: '400px',
                    borderRadius: '15px',
                    boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
                    transition: 'transform 0.3s ease, box-shadow 0.3s ease',
                }}
                hoverable
            >
                <Space 
                    direction="vertical" 
                    size="large" 
                    style={{ width: '100%' }}
                >
                    <Title 
                        level={1} 
                        style={{ 
                            textAlign: 'center',
                            margin: '0 0 20px 0',
                            color: '#1890ff'
                        }}
                    >
                        小微书
                    </Title>

                    <Button 
                        type="primary" 
                        icon={<UserOutlined />} 
                        size="large"
                        href="/users/login"
                        block
                        style={{ height: '45px' }}
                    >
                        登录/注册
                    </Button>

                    <Button 
                        type="default" 
                        icon={<ReadOutlined />}
                        size="large"
                        href="/articles/list"
                        block
                        style={{ height: '45px' }}
                    >
                        文章列表
                    </Button>

                    <Button 
                        type="default" 
                        icon={<EditOutlined />}
                        size="large"
                        href="/articles/edit"
                        block
                        style={{ height: '45px' }}
                    >
                        写文章
                    </Button>

                    <Space.Compact style={{ width: '100%' }}>
                        <Input 
                            placeholder="输入文章ID" 
                            value={articleId}
                            onChange={(e) => setArticleId(e.target.value)}
                            size="large"
                            style={{ 
                                borderTopLeftRadius: '6px',
                                borderBottomLeftRadius: '6px',
                            }}
                        />
                        <Button 
                            type="primary"
                            icon={<EyeOutlined />}
                            size="large"
                            href={`/articles/view?id=${articleId}`}
                            style={{ 
                                borderTopRightRadius: '6px',
                                borderBottomRightRadius: '6px',
                            }}
                        >
                            查看线上库文章
                        </Button>
                    </Space.Compact>
                </Space>
            </Card>
        </div>
    );
};

export default App;