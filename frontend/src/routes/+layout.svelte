<script>
	import { page } from '$app/stores';

	$: loggedIn = $page.data.user !== undefined;
</script>

<div>
	<nav>
		<ul>
			{#if loggedIn}
				<li>
					Hello, {$page.data.user.email}
				</li>
				<li>
					<a href="/">Home</a>
				</li>
			{:else}
				<li><strong>todo app</strong></li>
			{/if}
		</ul>
		<ul>
			{#if loggedIn}
				<li>
					<a href="/logout" data-sveltekit-preload-data="off" role="button">Logout</a>
				</li>
			{:else}
				<li>
					<a href="/login" role="button">Login</a>
				</li>
			{/if}
		</ul>
	</nav>
	<slot />
</div>

<style>
	nav {
		margin-bottom: 1em;
	}
	div {
		max-width: 30em;
		margin: 0 auto;
	}
</style>
