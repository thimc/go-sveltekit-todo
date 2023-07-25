<script>
	import { enhance } from '$app/forms';
	import { fly, fade } from 'svelte/transition';

	export let data;
</script>

<main class="centered">
	<form method="POST" action="?/create" use:enhance>
		<h1>my todos</h1>
		<label>
			<input type="text" autocomplete="off" name="description" />
		</label>
	</form>

	<ul class="todos">
		{#each data.todos as todo (todo.id)}
			<li in:fly={{ y: -120, duration: 120 }} out:fade={{duration: 200}}>
				<form method="POST" action="?/delete">
					<label>
						<input type="hidden" name="id" value={todo.id} />
						<input type="checkbox" checked={todo.done} />
						<span>{todo.description}</span>
						<button aria-label="Mark as complete">
							<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
								<path
									fill="#888"
									stroke="none"
									d="M22 4.2h-5.6L15 1.6c-.1-.2-.4-.4-.7-.4H9.6c-.2 0-.5.2-.6.4L7.6 4.2H2c-.4 0-.8.4-.8.8s.4.8.8.8h1.8V22c0 .4.3.8.8.8h15c.4 0 .8-.3.8-.8V5.8H22c.4 0 .8-.3.8-.8s-.4-.8-.8-.8zM10.8 16.5c0 .4-.3.8-.8.8s-.8-.3-.8-.8V10c0-.4.3-.8.8-.8s.8.3.8.8v6.5zm4 0c0 .4-.3.8-.8.8s-.8-.3-.8-.8V10c0-.4.3-.8.8-.8s.8.3.8.8v6.5z"
								/>
							</svg>
						</button>
					</label>
				</form>
			</li>
		{/each}
	</ul>
</main>

<style>
	:global(html) {
		background-color: #0c0c0c;
		color: #efefef;
	}

	.centered {
		max-width: 23em;
		margin: 0 auto;
	}

	h1 {
		text-align: center;
		font-size: 3em;
    margin-bottom: 0.3em;
	}

	input[type='text'] {
		width: 100%;
		margin-bottom: 2em;
		line-height: 2.2em;
		background-color: #222;
		border: 1px solid #444;
		border-radius: 0.5em;
		color: #efefef;
	}
	input[type='text']:focus {
		outline: 2px solid #666;
	}

	ul {
		padding: 0;
	}

	li {
		list-style: none;
		margin-bottom: 1em;
    border: 1px solid #111;
    border-radius: 10px;
    padding: 5px 5px;
    line-height: 1.6em;
	}

  label {
    display: flex;
  }

  span {
		width: 100%;
		margin-left: 1em;
		font-size: 1.1em;
	}

	button {
    margin-left: 0.2em;
		background-color: transparent;
		border: none;
		cursor: pointer;
		width: 36px;
		opacity: 0.5;
		transition: opacity 0.2s;
	}

	button:hover {
		opacity: 1;
	}
</style>
