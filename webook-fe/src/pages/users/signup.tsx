import React from 'react';
import { Button, Form, Input, Typography } from 'antd';
import { MailOutlined, LockOutlined } from '@ant-design/icons';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";

const { Title, Text } = Typography;

const onFinish = (values: any) => {
    axios.post("/users/signup", values)
        .then((res) => {
            if(res.status !== 200) {
                alert(res.statusText);
                return;
            }
            if(typeof res.data === 'string') {
                alert(res.data);
            } else {
                const msg = res.data?.msg || JSON.stringify(res.data);
                alert(msg);
                if(res.data.code === 0) {
                    router.push('/users/login');
                }
            }
        }).catch((err) => {
        alert(err);
    });
};

const onFinishFailed = () => {
    alert("输入有误");
};

const SignupForm = () => {
    return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh',
            backgroundColor: '#f0f2f5'
        }}>
            <div style={{
                width: '396px',
                padding: '20px',
                backgroundColor: '#fff',
                borderRadius: '8px',
                boxShadow: '0 2px 4px rgba(0, 0, 0, .1), 0 8px 16px rgba(0, 0, 0, .1)',
                textAlign: 'center'
            }}>
                <Title level={2} style={{ color: '#1877f2', marginBottom: '30px' }}>
                    注册
                </Title>
                <Form
                    name="signup"
                    initialValues={{ remember: true }}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="off"
                >
                    <Form.Item
                        name="email"
                        rules={[{ required: true, message: '请输入邮箱', type: 'email' }]}
                    >
                        <Input
                            prefix={<MailOutlined />}
                            size="large"
                            placeholder="邮箱"
                            style={{ borderRadius: '6px' }}
                        />
                    </Form.Item>

                    <Form.Item
                        name="password"
                        rules={[{ required: true, message: '请输入密码' }]}
                    >
                        <Input.Password
                            prefix={<LockOutlined />}
                            size="large"
                            placeholder="密码"
                            style={{ borderRadius: '6px' }}
                        />
                    </Form.Item>

                    <Form.Item
                        name="confirmPassword"
                        dependencies={['password']}
                        rules={[
                            { required: true, message: '请确认密码' },
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (!value || getFieldValue('password') === value) {
                                        return Promise.resolve();
                                    }
                                    return Promise.reject(new Error('两次输入的密码不一致'));
                                },
                            }),
                        ]}
                    >
                        <Input.Password
                            prefix={<LockOutlined />}
                            size="large"
                            placeholder="确认密码"
                            style={{ borderRadius: '6px' }}
                        />
                    </Form.Item>

                    <Form.Item>
                        <Button
                            type="primary"
                            htmlType="submit"
                            style={{
                                width: '65%',
                                height: '40px',
                                fontSize: '18px',
                                fontWeight: 'bold',
                                backgroundColor: '#1877f2',
                                borderColor: '#1877f2',
                                borderRadius: '6px'
                            }}
                        >
                            注册
                        </Button>
                    </Form.Item>
                </Form>

                <div style={{ borderBottom: '1px solid #dadde1', margin: '20px 0' }} />

                <Button
                    type="default"
                    style={{
                        width: '45%',
                        height: '40px',
                        fontSize: '15px',
                        fontWeight: 'bold',
                        backgroundColor: '#42b72a',
                        borderColor: '#42b72a',
                        color: '#fff',
                        marginBottom: '20px',
                        borderRadius: '6px'
                    }}
                >
                    <Link href="/users/login">返回登录</Link>
                </Button>
            </div>
        </div>
    );
};

export default SignupForm;