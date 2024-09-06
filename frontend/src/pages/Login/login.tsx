import React from 'react';
import { Form, Input, Button, Checkbox, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './login.css';

const Login: React.FC = () => {
    const navigate = useNavigate();

    const onFinish = async (values: any) => {
        try {
            const { data, status } = await axios.post('http://localhost:8000/api/signin', {
                email: values.username,
                password: values.password
            });

            if (status === 200) {
                const { email, id: memberID = 1, token = 'dummy-token' } = data;  // ใช้ค่า default เพื่อทดสอบ
                if (memberID && token) {
                    message.success('Login successful');
                    localStorage.setItem('isLogin', 'true');
                    localStorage.setItem('email', email); 
                    localStorage.setItem('memberID', memberID);
                    localStorage.setItem('token', token); 

                    const isAdmin = email === 'sa@gmail.com' && values.password === '123456';
                    localStorage.setItem('isAdmin', isAdmin ? 'true' : 'false');
                    navigate(isAdmin ? '/dashboard' : '/home');
                } else {
                    message.error('Invalid login response');
                }
            }
        } catch (err) {
            message.error('Login failed. Please check your username and password.');
        }
    };

    return (
        <div className="login-container">
            <div className="login-formContainer">
                <h1 className="login-title">MERJE CINIPLEX</h1>
                <Form
                    name="login"
                    onFinish={onFinish}
                    className="login-form"
                    layout="vertical"
                    requiredMark={false}
                >
                    <Form.Item
                        label={<span className="login-label">Username</span>}
                        name="username"
                        rules={[{ required: true, message: 'Please enter your username' }]}
                    >
                        <Input placeholder="Username" className="login-input" />
                    </Form.Item>

                    <Form.Item
                        label={<span className="login-label">Password</span>}
                        name="password"
                        rules={[{ required: true, message: 'Please enter your password' }]}
                    >
                        <Input.Password placeholder="Password" className="login-input" />
                    </Form.Item>

                    <Form.Item name="remember" valuePropName="checked" className="login-rememberMe">
                        <Checkbox className="login-checkbox">Remember me</Checkbox>
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" className="login-button">
                            Login
                        </Button>
                    </Form.Item>

                    <Form.Item className="login-registerContainer">
                        <span className="login-registerText">Don't have an account?</span>
                        <Button type="default" className="login-registerButton" onClick={() => navigate('/signup')}>
                            Register
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        </div>
    );
};

export default Login;
