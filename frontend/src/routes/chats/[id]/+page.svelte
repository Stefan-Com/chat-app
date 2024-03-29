<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import type { ChatRoom, Message, User } from '$lib/interfaces.ts';
	import { GetRoom, FetchMessages, GetUserData } from '$lib/chat-rooms.ts';

	const id: string = $page.params.id;
	let currentRoom: ChatRoom = {};

	let messages: Message[] = [];
	let socket: WebSocket;
	let msg = '';

	let USER: User | any = {};

	onMount(async () => {
		socket = new WebSocket(`ws://localhost:7000/api/socket/${id}`);
		socket.onopen = () => {
			socket.onmessage = (event) => {
				messages = [...messages, JSON.parse(event.data)];
			};
		};
		currentRoom = await GetRoom(id);
		messages = await FetchMessages(currentRoom.messages);
		if (!messages) messages = [];
		USER = await GetUserData();
		if (!USER) USER = {};
	});

	function formatDate(dateStr: string) {
		if (!dateStr) return;
		const date = new Date(dateStr);
		const hour = date.getHours() > 12 ? date.getHours() - 12 : date.getHours();
		const minute = date.getMinutes() > 9 ? date.getMinutes() : '0' + date.getMinutes();
		const meridian = date.getHours() > 12 ? 'PM' : 'AM';

		return `${hour}:${minute} ${meridian}`;
	}
</script>

<div class="container">
	<div>
		<h2>Room {currentRoom.title}</h2>
		<div class="msg-container">
			{#each messages as message (message.id)}
				<div class="msg-content" class:sent-by-me={USER.id === message.fromid}>
					<div>
						{message.text}
					</div>
					<span>{formatDate(String(message.sentat))}</span>
				</div>
			{/each}
		</div>
	</div>

	<div class="msg-input">
		<textarea bind:value={msg} />
		<button
			type="button"
			on:click={() => {
				if (!msg.length) return;
				const data = JSON.stringify({ text: msg, toid: parseInt(id) });
				socket.send(data);

				msg = '';
			}}>Send</button
		>
	</div>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
	}

	.container .msg-input textarea {
		resize: none;
		width: 674px;
		height: 29px;
		font-size: 15px;
		font-family: inherit;
		border-radius: 25px;
		padding: 10px;
	}

	.container .msg-input {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 20px;
		margin-top: auto;
		position: fixed;
		bottom: 24px;
		left: 0;
		right: 0;
	}

	.container .msg-input button {
		background: hsl(0, 0%, 100%);
		border: #005fcf solid 2px;
		cursor: pointer;
		border-radius: 25px;
		width: 97px;
		height: 49px;
		color: #005fcf;
		font-size: 15px;
		transition:
			background 0.5s,
			color 0.5s;
	}

	.container .msg-input button:focus,
	.container .msg-input button:hover {
		background: #005fcf;
		border-color: #fff;
		color: #fff;
	}

	.msg-content {
		width: fit-content;
		height: fit-content;
		min-width: 250px;
		padding: 26px 14px 12px 14px;
		border-radius: 20px;
		word-wrap: break-word;
		color: #000;
		box-shadow: 1px 3px 3px 4px rgba(0, 0, 0, 0.25);
		background: #fff;
	}

	.msg-content div {
		font-size: 18px;
		max-width: 350px;
	}

	.msg-content span {
		display: flex;
		justify-content: flex-end;
		color: #c5c5c5;
		font-size: 12px;
		margin-top: 4px;
	}

	.msg-container {
		display: flex;
		flex-direction: column;
		gap: 21px;
		max-height: 535px;
		width: 100%;
		overflow-y: scroll;
		background-position-y: bottom;
	}

	.sent-by-me {
		background: #479cff;
		color: #fff;
		border: none;
		margin-left: auto;
	}
</style>
