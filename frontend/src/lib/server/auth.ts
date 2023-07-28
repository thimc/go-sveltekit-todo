import { API_URL } from '$env/static/private';
import type { RequestEvent } from '@sveltejs/kit';

export const authenticateUser = async (event: RequestEvent): Promise<App.User | null> => {
	const token = event.cookies.get('jwt');
	if (token === undefined) return null;

	try {
		const res = await fetch(`${API_URL}/api/check`, {
			method: 'GET',
			headers: { Authorization: `Bearer ${token}` }
		});

		const result = await res.json();

		if (result.success === false) {
			event.cookies.delete('jwt');
			return null;
		}

		const user: App.User = {
			id: result.id,
			email: result.email,
			token: token
		};
		return user;
	} catch (err) {
		console.log('authenticateUser error:', err);
	}
	return null;
};
