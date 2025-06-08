import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 100,
    duration: '30s',
    thresholds: {
        http_reqs: ['rate>100'],
        http_req_duration: ['p(95)<1000'],
    },
};

const BASE_URL = 'http://localhost:9000/external/api/v1';
const USER = { username: 'admin', password: 'adminpass123' };

export default function () {
    // 1. Login to get token
    let loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify(USER), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(loginRes, { 'login success': (r) => r.status === 200 });
    const token = loginRes.json('data.token') || loginRes.json('data.access_token');
    const authHeaders = {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    };

    const payrollId = "c7c50c3b-5d94-4aa1-b9df-17b5f184b5bb"

    // 4. List payslips with random pagination
    const totalPayslips = 100;
    const pageSize = Math.floor(Math.random() * 10) + 1; // 1-10
    const maxPage = Math.ceil(totalPayslips / pageSize);
    const page = Math.floor(Math.random() * maxPage) + 1;

    // Simulate users hitting all endpoints together (parallel requests)
    const requests = [
        [
            'POST', `${BASE_URL}/attendance`, JSON.stringify({}), authHeaders
        ],
        [
            'GET', `${BASE_URL}/payroll/${payrollId}/payslip`, null, authHeaders
        ],
        [
            'GET', `${BASE_URL}/payroll/${payrollId}/payslips?limit=${pageSize}&page=${page}`, null, authHeaders
        ],
        [
            'POST', `${BASE_URL}/reimbursement`, JSON.stringify({
                "date": "2025-07-08",
                "amount": 2000000,
                "description": "Hotels Ticket"
            }), authHeaders
        ]
    ];
    const responses = http.batch(requests);
    check(responses[0], { 'attendance success': (r) => r.status === 200 });
    check(responses[1], { 'payslip success': (r) => r.status === 200 });
    check(responses[2], { 'payslips list success': (r) => r.status === 200 });
    check(responses[3], { 'reimbursement success': (r) => r.status === 200 });
    sleep(1);
}
