import React from 'react';
import { Button, Form, Input } from 'antd';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";

const onFinish = (values: any) => {
    axios.post("/users/login", values)
        .then((res) => {
            if(res.status != 200) {
                alert(res.statusText);
                return
            }
            if(typeof res.data == 'string') {
                alert(res.data);
            } else {
                const msg = res.data?.msg || JSON.stringify(res.data)
                alert(msg);
                if(res.data.code == 0) {
                    router.push('/articles/list')
                }
            }
        }).catch((err) => {
        alert(err);
    })
};

const onFinishFailed = (errorInfo: any) => {
    alert("输入有误")
};

// const LoginForm: React.FC = () => {
//     return (<Form
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
//         <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
//             <Button type="primary" htmlType="submit">
//                 登录
//             </Button>
//             <Link href={"/users/login_sms"} >
//                 &nbsp;&nbsp;手机号登录
//             </Link>
//             <Link href={"/users/login_wechat"} >
//                 &nbsp;&nbsp;微信扫码登录
//             </Link>
//             <Link href={"/users/signup"} >
//                 &nbsp;&nbsp;注册
//             </Link>
//         </Form.Item>
//     </Form>
// )};

//代码升级版本
const LoginForm: React.FC = () => {
    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', backgroundColor: '#f0f2f5' }}>
            <Form
                name="basic"
                style={{ maxWidth: 1000, background: '#fff', padding: '60px', borderRadius: '8px', boxShadow: '0 2px 10px rgba(0, 0, 0, 0.1)' }}
                initialValues={{ remember: true }}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                autoComplete="off"
            >
                <h2 style={{ textAlign: 'center', marginBottom: '30px', color: '#1890ff', fontSize: '24px' }}>登录</h2>

                <Form.Item name="email" rules={[{ required: true, message: '请输入邮箱' }]}>
                    <Input size="large" placeholder="邮箱" style={{ fontSize: '18px', width: '100%', boxShadow: '0 1px 5px rgba(0, 0, 0, 0.2)' }} />
                </Form.Item>

                <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                    <Input.Password size="large" placeholder="密码" style={{ fontSize: '18px', width: '100%', boxShadow: '0 1px 5px rgba(0, 0, 0, 0.2)' }} />
                </Form.Item>

                <Form.Item wrapperCol={{ offset: 0, span: 24 }} style={{ textAlign: 'center' }}>
                    <Button type="primary" htmlType="submit" style={{ backgroundColor: '#1890ff', borderColor: '#1890ff', width: '100%', fontSize: '20px', padding: '10px' }}>
                        登录
                    </Button>
                    <div style={{ marginTop: '20px' }}>
                        <Link href={"/users/login_sms"} style={{ fontSize: '18px' }}>手机号登录</Link>
                        <Link href={"/users/login_wechat"} style={{ marginLeft: '10px', fontSize: '18px' }}>微信扫码登录</Link>
                    </div>
                    <div style={{ marginTop: '10px' }}>
                        <Link href={"/users/signup"} style={{ fontSize: '18px' }}>注册</Link>
                    </div>
                </Form.Item>
            </Form>
        </div>
    );
};

export default LoginForm;