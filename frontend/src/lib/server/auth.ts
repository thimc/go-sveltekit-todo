import type { RequestEvent } from '@sveltejs/kit';

export const authenticateUser = async (event: RequestEvent): Promise<User | null> => {
	const token = event.cookies.get('jwt');
	if (token === undefined) return null;

	const res = await fetch(`http://localhost:1234/api/check`, {
		method: 'GET',
		headers: { Authorization: `Bearer ${token}` }
	});

	const result = await res.json();

	if (result.success === false) {
		event.cookies.delete('jwt');
		return null;
	}

	const user: User = {
		id: result.id,
		email: result.email,
		token: token
	};
	return user;
};