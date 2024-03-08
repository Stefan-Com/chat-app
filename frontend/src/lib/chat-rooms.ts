import type { ChatRoom } from "$lib/interfaces.ts"

export async function AddChatRoom(info: ChatRoom) {
		const response = await fetch("http://localhost:7000/rooms", {
			method: "POST",
			body: JSON.stringify(info),
			credentials: "include",
			headers: {
				"Content-Type": "application/json"
			}
		}).then(res => res.json())

		console.log(response)
		if (response.status === "success") return response.response
		return {}
}

export async function GetChatRooms() {
		const response = await fetch('http://localhost:7000/rooms', {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		}).then((res) => res.json());
		if (response.status === "success") return response.response
		return []
}

export async function GetRoom(id: string) {
		const response = await fetch(`http://localhost:7000/rooms/${id}`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		}).then((res) => res.json());
		if (response.status === 'success') return response.response;
		return {};
}

export async function GetUsers() {
	const response = await fetch(`http://localhost:7000/users`, {
		method: 'GET',
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json'
		}
	}).then((res) => res.json());
	if (response.status === 'success') return response.response;
	return []
}

export async function FetchMessages(id: number[]) {
		const body = JSON.stringify(id);
		const data = await fetch('http://localhost:7000/messages', {
			method: 'PUT',
			body: body,
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		}).then((res) => res.json());
		if (data.status === 'success') return data.response;
		return [];
}

