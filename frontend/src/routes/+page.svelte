<script>
  import { enhance } from '$app/forms';
	import { fly, fade } from 'svelte/transition';

	export let data;
	export let form;
</script>

<div>
	<form method="POST" action="?/createTodo" use:enhance>
		<!--		<label for="todo">my todos</label>-->
		<input type="text" name="content" placeholder="Create a new todo..." />

		{#if form?.success === false}
			<p class="error">Error: {form?.message}</p>
		{/if}
	</form>

<!--
	<pre>{JSON.stringify(data.user, null, 2)}</pre>
	<pre>{JSON.stringify(data.todos, null, 2)}</pre>
-->

	<div class="todos">
		{#each data.todos as todo (todo.id)}
			<div class="todoItem" in:fly={{ y: -120, duration: 120 }} out:fade={{ duration: 200 }}>
				<div class="todoControls">
					<input type="checkbox" name="done" bind:value={todo.done} />
				</div>
				<div
					class="todoContent"
					style={todo.done ? 'text-decoration:line-through' : 'text-decoration: none'}
				>
					<!-- <h4 style="margin-bottom: 0">{todo.title}</h4> -->
					<span style="margin-bottom: 0">{todo.content}</span>
				</div>
				<form method="POST" action="?/deleteTodo" class="todoControls" use:enhance>
					<input type="hidden" value={todo.id} name="id" />
					<button class="green">
						<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
							<path
								fill="currentColor"
								stroke="none"
								d="M22 4.2h-5.6L15 1.6c-.1-.2-.4-.4-.7-.4H9.6c-.2 0-.5.2-.6.4L7.6 4.2H2c-.4 0-.8.4-.8.8s.4.8.8.8h1.8V22c0 .4.3.8.8.8h15c.4 0 .8-.3.8-.8V5.8H22c.4 0 .8-.3.8-.8s-.4-.8-.8-.8zM10.8 16.5c0 .4-.3.8-.8.8s-.8-.3-.8-.8V10c0-.4.3-.8.8-.8s.8.3.8.8v6.5zm4 0c0 .4-.3.8-.8.8s-.8-.3-.8-.8V10c0-.4.3-.8.8-.8s.8.3.8.8v6.5z"
							/>
						</svg>
					</button>
				</form>
			</div>
		{/each}
	</div>
</div>

<style>
	label {
		font-size: 32px;
		text-align: center;
	}
	div {
		margin-bottom: 15px;
	}
	.todos > * {
		padding: calc(var(--spacing) / 2) 0;
		border-radius: var(--border-radius);
		background: var(--code-background-color);
		text-align: center;
	}

	.todoItem {
		display: flex;
	}
	.todoContent {
    margin: auto;
		width: 100%;
	}
	.todoControls {
    margin: auto;
		display: flex;
		margin-left: 10px;
		margin-right: 10px;
	}

	.todoControls > button {
    margin: auto;
		background-color: inherit;
		padding: 0;
		border: none;
		cursor: pointer;
		opacity: 0.5;
		transition: opacity 0.2s;
		justify-items: center;
	}

	.todoControls > button:hover {
		opacity: 1;
	}

	.line-through {
		text-decoration: line-through;
	}

	.error {
		color: var(--del-color);
		font-weight: bold;
	}
</style>
