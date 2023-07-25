import * as db from '$lib/server/database';
import { user } from '$lib/userStore';

/** @type {import('./$types').PageServerLoad} */
export function load({}) {
	let userId = '';
	user.subscribe((id) => {
		userId = id;
	});

	if (!userId) {
		// cookies.set('userId', crypto.randomUUID());
		user.set(crypto.randomUUID());
	}
	return {
		todos: db.getTodos(userId) ?? []
	};
}

/** @type {import('./$types').Actions} */
export const actions = {
	create: async ({ request }) => {
		await new Promise((fulfil) => setTimeout(fulfil, 200));

		const data = await request.formData();
		const description = String(data.get('description'));
		let userId = '';
		user.subscribe((id) => {
			userId = id;
		});
		// const userId = cookies.get('userId') ?? '';

		db.createTodo(userId, description);
	},
	delete: async ({ request }) => {
		const data = await request.formData();
		const todoId = String(data.get('id'));
		let userId = '';
		user.subscribe((id) => {
			userId = id;
		});
		db.deleteTodo(userId, todoId);
	},
	update: async ({ request }) => {
		const data = await request.formData();
		const todoId = String(data.get('todoId'));
		const todoDone = String(data.get('todoDone')) == 'true' ? true : false;

		let userId = '';
		user.subscribe((id) => {
			userId = id;
		});

		db.updateTodo(userId, todoId, todoDone);
	}
};
