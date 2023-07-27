import { authenticateUser } from '$lib/server/auth';
import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ resolve, event }) => {
	event.locals.user = await authenticateUser(event);

	if (event.url.pathname == '/') {
		if (!event.locals.user) {
			throw redirect(303, '/login');
		}
	}

	return await resolve(event);
};
