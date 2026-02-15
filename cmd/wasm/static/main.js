const go = new Go();

window.cellSize = getCellSize();

WebAssembly.instantiateStreaming(
	fetch("main.wasm"), 
	go.importObject).then((result) => {
		go.mem = result.instance.exports.mem;
    	go.run(result.instance);
	});

function updateCanvas(ptr, width, height, status) {
	let b = go.mem.buffer;
	let len = width * height * 4;
	let pixels = new Uint8ClampedArray(b, ptr, len);
	let imageData = new ImageData(pixels, width, height);
	writeCtx.putImageData(imageData, 0, 0);
	ctx.drawImage(hiddenCanvas, 0, 0, innerWidth, innerHeight);
	ctx.fillText(status, 10, 30);
}

function getCellSize() {
	let search = window.location.search;
	let params = new URLSearchParams(search);
	let sizeParam = params.get("cellSize") ?? 3;
	return +sizeParam | 0;
}

// removed closure untill flow is repaired.
	let canvas = window.document.getElementById("canvas");
	let innerWidth = window.innerWidth;
	let innerHeight = window.innerHeight;
	canvas.width = innerWidth;
	canvas.height = innerHeight;
	let ctx = canvas.getContext("2d");
	
	ctx.font = "40px monospace";
	ctx.fillStyle = "red";
	ctx.imageSmoothingEnabled = false;

	let width = innerWidth / window.cellSize;
	let height = innerHeight / window.cellSize;

	let hiddenCanvas = window.document.createElement("canvas");
	hiddenCanvas.width = width;
	hiddenCanvas.height = height;
	hiddenCanvas.id = "hidden";

	let writeCtx = hiddenCanvas.getContext("2d");

	console.log("init end", typeof(cellSize), cellSize);
