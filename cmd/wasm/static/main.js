import { RollingAverage } from './stats.js'
const go = new Go();

window.gol = {
	status: [""],
	lastStatusUpdate: 0,
	dims: {
		cellSize: getCellSize(),
		width: 0,
		height: 0
	},
	buffer: null,
	t: {
		frame: new RollingAverage(60),
		tick: new RollingAverage(60),
		paint: new RollingAverage(60),
		next: new RollingAverage(60),
		render: new RollingAverage(60)
	}
}

function updateCanvas() {
	const gol = window.gol;

	gol.writeCtx.putImageData(gol.imageData, 0, 0);
	gol.ctx.drawImage(gol.hiddenCanvas, 0, 0, 
		gol.canvas.width, 
		gol.canvas.height);
	
	drawStatusBlock(gol.ctx, gol.status, 16, 64, 45);
}
function drawStatusBlock(ctx, lines, x, y, lineHeight) {
    ctx.save();
    ctx.textBaseline = "top"; 
    ctx.textAlign = "start";
    
    lines.forEach((line, i) => {
        ctx.fillText(line, x, y + (i * lineHeight));
    });
    ctx.restore();
}


function updateStatus(status) {
	window.gol.status = status;
}

function addGoTimings(next, render) {
	window.gol.t.next.add(next);
	window.gol.t.render.add(render);
}

function getCellSize() {
	let search = window.location.search;
	let params = new URLSearchParams(search);
	let sizeParam = params.get("cellSize") ?? 3;
	return int(sizeParam);
}
function int(number) {
	return +number | 0;
}

function renderFrame(ts) {
	performance.mark("start-frame");
	tick();
	performance.mark("update-canvas");
	if (ts - window.gol.lastStatusUpdate > 500) {
		updateStatus([
			`Frame: ${window.gol.t.frame.average().toFixed(2)}ms`,
			`├─Tick: ${window.gol.t.tick.average().toFixed(2)}ms`,
			`│ ├─Next:   ${window.gol.t.next.average().toFixed(0)}µs`,
			`│ └─Render: ${window.gol.t.render.average().toFixed(0)}µs`,
			`└─Paint: ${window.gol.t.paint.average().toFixed(2)}ms`
		]);
		
		window.gol.lastStatusUpdate = ts;
	}
	if (go.mem.buffer.byteLength == 0) {
		setMemoryView(window.gol.ptr);
	}
	updateCanvas()

	performance.mark("end-frame");

	performance.measure("tick", "start-frame", "update-canvas");
	performance.measure("paint", "update-canvas", "end-frame");
	performance.measure("frame", "start-frame", "end-frame");

	return requestAnimationFrame(renderFrame);
}

function setMemoryView(ptr) {
	const gol = window.gol;
	const dims = gol.dims;
	const len = dims.width * dims.height * 4;
	
	const memoryView = new Uint8ClampedArray(
		go.mem.buffer, ptr, len
	);

	const imageData = new ImageData(
		memoryView, dims.width, dims.height
	);

	gol.ptr = ptr;
	gol.imageData = imageData;
}

function updateDimensions(canvas, hiddenCanvas) {
	canvas.width = window.innerWidth;
	canvas.height = window.innerHeight;

	let cs = getCellSize();

	hiddenCanvas.width = int(canvas.width / cs)
	hiddenCanvas.height = int(canvas.height / cs)	

	return {
		cellSize: cs,
		width: hiddenCanvas.width,
		height: hiddenCanvas.height
	}
}


const observer = new PerformanceObserver((list) => {
	const entries = list.getEntries();

	for (const entry of entries) {
		window.gol.t[entry.name].add(entry.duration);
	}
});

observer.observe({ 'entryTypes': ['measure'] });

async function init() {
	const gol = window.gol;
    // 1. Get DOM elements (canvas, hiddenCanvas)
	let canvas = window.document.getElementById("canvas");
	gol.canvas = canvas;
	let hiddenCanvas = window.document.createElement("canvas");
	hiddenCanvas.id = "hidden";
	gol.hiddenCanvas = hiddenCanvas;
	gol.dims = updateDimensions(canvas, hiddenCanvas);
    // 2. Setup Contexts (ctx, writeCtx)
    
	let ctx = canvas.getContext("2d");
	ctx.font = "48px 'Roboto Mono', 'Source Code Pro', Consolas, monospace";
	ctx.fillStyle = "red";
	ctx.textRendering = "geometricPrecision";
	ctx.imageSmoothingEnabled = false;
	gol.ctx = ctx;
	
	gol.writeCtx = hiddenCanvas.getContext("2d");

	// Exposed functions
	const fn = gol.fn = {};
	fn.setMemoryView = setMemoryView;
    fn.renderFrame = renderFrame;
    fn.addGoTimings = addGoTimings;

    const result = await WebAssembly.instantiateStreaming(
        fetch("main.wasm"), go.importObject);
        
    go.mem = result.instance.exports.mem;

    go.run(result.instance);
}

init();


