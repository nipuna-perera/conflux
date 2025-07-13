// API proxy to forward requests to the backend
// This handles cases where the frontend needs to proxy API calls
import { json, error } from '@sveltejs/kit';

const BACKEND_URL = 'http://backend:8080/api';

export async function GET({ params, url, request }: { params: any, url: URL, request: Request }) {
	return proxyRequest(request, params.path, url);
}

export async function POST({ params, url, request }: { params: any, url: URL, request: Request }) {
	return proxyRequest(request, params.path, url);
}

export async function PUT({ params, url, request }: { params: any, url: URL, request: Request }) {
	return proxyRequest(request, params.path, url);
}

export async function DELETE({ params, url, request }: { params: any, url: URL, request: Request }) {
	return proxyRequest(request, params.path, url);
}

async function proxyRequest(request: Request, path: string, url: URL) {
	try {
		// Forward the request to the backend
		const backendUrl = `${BACKEND_URL}/${path}${url.search}`;
		
		const headers: Record<string, string> = {};
		
		// Copy relevant headers
		for (const [key, value] of request.headers.entries()) {
			if (['authorization', 'content-type', 'accept'].includes(key.toLowerCase())) {
				headers[key] = value;
			}
		}
		
		const body = ['GET', 'HEAD'].includes(request.method) ? undefined : await request.text();
		
		const response = await fetch(backendUrl, {
			method: request.method,
			headers,
			body,
		});
		
		const responseText = await response.text();
		
		// Try to parse as JSON, fallback to text
		let responseData;
		try {
			responseData = JSON.parse(responseText);
		} catch {
			responseData = responseText;
		}
		
		if (!response.ok) {
			throw error(response.status, responseData);
		}
		
		return json(responseData);
	} catch (err) {
		console.error('API proxy error:', err);
		throw error(500, 'Internal server error');
	}
}
