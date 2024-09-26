// import React from 'react';
// import { Button, Form, Input } from 'antd';
// import axios from "@/axios/axios";
// import Link from "next/link";
// import router from "next/router";
//
// const onFinish = (values: any) => {
//     axios.post("/users/signup", values)
//         .then((res) => {
//             if(res.status != 200) {
//                 alert(res.statusText);
//                 return
//             }
//             if(typeof res.data == 'string') {
//                 alert(res.data);
//             } else {
//                 const msg = res.data?.msg || JSON.stringify(res.data)
//                 alert(msg);
//                 if(res.data.code == 0) {
//                     router.push('/users/login')
//                 }
//             }
//
//         }).catch((err) => {
//             alert(err);
//     })
// };
//
// const onFinishFailed = (errorInfo: any) => {
//     alert("输入有误")
// };
//
// const SignupForm: React.FC = () => (
//     <Form
//         name="basic"
//         labelCol={{ span: 8 }}
//         wrapperCol={{ span: 16 }}
//         style={{ maxWidth: 600 }}
//         initialValues={{ remember: true }}
//         onFinish={onFinish}
//         onFinishFailed={onFinishFailed}
//         autoComplete="off"
//     >
//         <Form.Item
//             label="邮箱"
//             name="email"
//             rules={[{ required: true, message: '请输入邮箱' }]}
//         >
//             <Input />
//         </Form.Item>
//
//         <Form.Item
//             label="密码"
//             name="password"
//             rules={[{ required: true, message: '请输入密码' }]}
//         >
//             <Input.Password />
//         </Form.Item>
//
//         <Form.Item
//             label="确认密码"
//             name="confirmPassword"
//             rules={[{ required: true, message: '请确认密码' }]}
//         >
//             <Input.Password />
//         </Form.Item>
//         <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
//             <Button type="primary" htmlType="submit">
//                 注册
//             </Button>
//             <Link href={"/users/login"}>&nbsp;登录</Link>
//         </Form.Item>
//     </Form>
// );
//
// export default SignupForm;


//升级版本

import React from 'react';
import { Button, Form, Input, Typography, ConfigProvider } from 'antd';
import { MailOutlined, LockOutlined } from '@ant-design/icons';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";

const { Title, Text } = Typography;

// 更柔和的蓝色主题
const theme = {
    token: {
        colorPrimary: '#1877f2',
        colorLink: '#1877f2',
        borderRadius: 8,
    },
};

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

const SignupForm = () => (
    <ConfigProvider theme={theme}>
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh',
            background: '#f0f2f5'
        }}>
            <div style={{
                width: '100%',
                maxWidth: '400px',
                padding: '20px',
                textAlign: 'center'
            }}>
                <Title level={1} style={{ color: '#1877f2', marginBottom: '30px' }}>小微书</Title>
                <div style={{
                    background: 'white',
                    padding: '20px',
                    borderRadius: '8px',
                    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1), 0 8px 16px rgba(0, 0, 0, 0.1)'
                }}>
                    <Form
                        name="signup"
                        initialValues={{ remember: true }}
                        onFinish={onFinish}
                        onFinishFailed={onFinishFailed}
                        layout="vertical"
                        size="large"
                    >
                        <Form.Item
                            name="email"
                            rules={[{ required: true, message: '请输入邮箱', type: 'email' }]}
                        >
                            <Input prefix={<MailOutlined />} placeholder="邮箱" />
                        </Form.Item>

                        <Form.Item
                            name="password"
                            rules={[{ required: true, message: '请输入密码' }]}
                        >
                            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
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
                            <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
                        </Form.Item>

                        <Form.Item>
                            <Button type="primary" htmlType="submit" block style={{ height: '48px', fontSize: '16px', fontWeight: 'bold' }}>
                                注册
                            </Button>
                        </Form.Item>
                    </Form>
                </div>
                <div style={{ marginTop: '20px' }}>
                    <Text>已有账号？</Text> <Link href="/users/login" style={{ color: '#1877f2', fontWeight: 'bold' }}>立即登录</Link>
                </div>
            </div>
        </div>
    </ConfigProvider>
);

export default SignupForm;