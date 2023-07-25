/**
 * @typedef {Object} todoItem
 * @property {string} id
 * @property {string} description
 * @property {number} created
 * @property {number|null} updated
 * @property {boolean} done
 */

/** @type {Map<string, Array<todoItem>>} */
const db = new Map();

/**
 * Returns todos that matches the `userId`
 * @param {string} userId
 * @returns {todoItem[]|undefined} the todos
 */
export function getTodos(userId) {
	if (!db.get(userId)) {
		db.set(userId, [
			{
				id: crypto.randomUUID(),
				description: 'Learn SvelteKit and make awesome websites! 👍',
				created: Number(new Date()),
				updated: null,
				done: false
			}
		]);
	}

	return db.get(userId);
}

/**
 * Adds a new todo to the database
 * @param {string} userId
 * @param {string} description
 */
export function createTodo(userId, description) {
	const todos = db.get(userId);
	if (todos === undefined) {
		console.log('create: unknown user ID');
		return;
	}
	todos.push({
		id: crypto.randomUUID(),
		description,
		created: Number(new Date()),
		updated: null,
		done: false
	});
}

/**
 * Removes a todo by the `todoId`
 * @param {string} userId
 * @param {string} todoId
 */
export function deleteTodo(userId, todoId) {
	const todos = db.get(userId);
	if (todos === undefined) {
		console.log('delete: unknown user ID');
		return;
	}
	const index = todos.findIndex((todo) => {
		return todo.id === todoId;
	});

	if (index !== -1) todos.splice(index, 1);
}

/**
 * Updates the todo `done` status by the `todoId`. Requires an `userId`
 * @param {string} userId
 * @param {string} todoId
 * @param {boolean} done
 */
export function updateTodo(userId, todoId, done) {
	const todos = db.get(userId);
	if (todos === undefined) {
		console.log('update: unknown user ID');
		return;
	}

	todos.forEach((e) => {
		if (e.id === todoId) {
			e.updated = Number(new Date());
			e.done = !done;
		}
	});
}
