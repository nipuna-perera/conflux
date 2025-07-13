// API proxy for login endpoint
import { json, error } from '@sveltejs/kit';

const BACKEND_URL = 'http://backend:8080/api';

export async function POST({ request }: { request: Request }) {
	try {
		console.log('Login proxy called');
		
		const headers: Record<string, string> = {};
		
		// Copy relevant headers
		for (const [key, value] of request.headers.entries()) {
			if (['authorization', 'content-type', 'accept'].includes(key.toLowerCase())) {
				headers[key] = value;
			}
		}
		
		const body = await request.text();
		console.log('Forwarding to:', `${BACKEND_URL}/auth/login`);
		console.log('Body:', body);
		
		const response = await fetch(`${BACKEND_URL}/auth/login`, {
			method: 'POST',
			headers,
			body,
		});
		
		const responseText = await response.text();
		console.log('Backend response:', responseText);
		
		// Try to parse as JSON, fallback to text
		let responseData;
		try {
			responseData = JSON.parse(responseText);
		} catch {
			responseData = responseText;
		}
		
		if (!response.ok) {
			console.error('Backend error:', response.status, responseData);
			throw error(response.status, responseData);
		}
		
		return json(responseData);
	} catch (err) {
		console.error('API proxy error:', err);
		throw error(500, 'Internal server error');
	}
}
