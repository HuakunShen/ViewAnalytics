<!-- <script lang="ts">
	import { page } from '$app/stores';
	import { trpc } from '$lib/trpc/client';
  
	let greeting = 'press the button to load data';
	let loading = false;
  
	const loadData = async () => {
	  loading = true;
	  greeting = await trpc($page).greeting.query();
	  loading = false;
	};
  </script>
  
  <h6>Loading data in<br /><code>+page.svelte</code></h6>
  
  <a
	href="#load"
	role="button"
	class="secondary"
	aria-busy={loading}
	on:click|preventDefault={loadData}>Load</a
  >
  <p>{greeting}</p>
   -->
<script lang="ts">
	import { dev } from '$app/environment';
	import { Button } from '$lib/components/ui/button';
	import PocketBase from 'pocketbase';
	import { onMount } from 'svelte';

	const pb = new PocketBase(dev ? 'http://localhost:8090' : 'https://proxy.huakun.tech');

	onMount(async () => {
		const records = await pb.collection('proxy_records').getFullList();
		console.log(records);
	});

	async function onSignInClicked() {
		pb.collection('users')
			.authWithOAuth2({ provider: 'github' })
			.then((authData) => {
				console.log(authData);
				console.log(pb.authStore.isValid);
				console.log(pb.authStore.token);
				console.log(pb.authStore.model?.id);
			})
			.catch((err) => {
				console.error(err);
			});
	}
	function logout() {
		pb.authStore.clear();
	}
</script>

<Button onclick={onSignInClicked}>Sign In</Button>
<Button onclick={logout}>Logout</Button>
<span>Is Valid: {pb.authStore.isValid}</span>
<pre>Model: {JSON.stringify(pb.authStore.model, null, 2)}</pre>
