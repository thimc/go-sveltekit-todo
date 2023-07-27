import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	return {
		cookie: cookies.get('jwt')
	};
};

export const actions: Actions = {
	login: async ({ request, cookies }) => {
		const formData = await request.formData();
		const email = formData.get('email');
		const password = formData.get('password');

		try {
			const res = await fetch(`http://localhost:1234/api/login`, {
				method: 'POST',
				body: JSON.stringify({ email, password })
			});
			const result = await res.json();
			if (result.success == false) {
				console.log('Login failed result:', result);
				return fail(405, {
					success: false,
					message: result.message,
					email
				});
			}
			console.log('Login success result:', result);
			const now = Math.floor(Number(new Date()) / 1000);
			console.log('Expire:', result.expiresAt - now);
			cookies.set('jwt', result.token, {
				path: '/',
				sameSite: 'strict',
        httpOnly: true,
				secure: false,
				maxAge: result.expiresAt - now
			});
		} catch (err) {
			console.log('Error when connecting to the API', err);
			return fail(400, {
				success: false,
				message: `the API is not responding`,
				email
			});
		}
		throw redirect(300, '/');
	}
};
