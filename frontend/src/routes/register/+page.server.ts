import { fail, redirect } from '@sveltejs/kit';

export const actions: Actions = {
	default: async ({ request }) => {
		const formData = await request.formData();
		const email = formData.get('email');
		const password = formData.get('password');
		const passwordConfirm = formData.get('passwordConfirm');

		if (password !== passwordConfirm) {
			return fail(405, {
				success: false,
				message: 'the passwords does not match!',
				email
			});
		}

		try {
			const res = await fetch(`http://localhost:1234/api/register`, {
				method: 'POST',
				body: JSON.stringify({ email, password })
			});
			const result = await res.json();
			if (result.success == false) {
				console.log('Register failed result:', result);
				return fail(405, {
					success: false,
					message: result.message,
					email
				});
			}
      console.log('Register result:', result);
		} catch (err) {
			console.log('Register error:', err);
		}

    throw redirect(301, `/login?registeredEmail=${email}`);
	}
};
