export function syncParamsToUrl(params) {
	const url = new URL(window.location);
	Object.keys(params).forEach(key => {
		url.searchParams.set(key, params[key])
	});
	window.history.replaceState({}, '', url)
}
