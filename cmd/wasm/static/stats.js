export class RollingAverage {
    #sampleSize;
    #events;
    #count = 0;
    #index = 0;
    #sum = 0;

    constructor(sampleSize = 60) {
        this.#sampleSize = sampleSize;
        this.#events = new Float32Array(sampleSize);

        this.average = () => {
        	if (this.#count === 0) return 0;
        	if (this.#count < this.#sampleSize) {
        		return this.#sum / this.#count;
        	}
        	this.average = this._average.bind(this);
        	return this.average()
        }
    }

    add(event) {
    	this.#sum -= this.#events[this.#index];
        this.#events[this.#index] = event;
        this.#sum += event;
        this.#index++;
        if (this.#index >= this.#sampleSize) {
        	this.#index = 0;
        }
        this.#count++;
    }

    _average() {
        return this.#sum / this.#sampleSize;
    }
}

export class ExponentialAverage {
	#alpha;
	#current;

	constructor(alpha = 0.05) {
		this.#alpha = alpha;
		this.#current = null;

		this.add = (value) => {
			this.#current = value;
			this.add = this._add.bind(this);
		}
	}

	_add(value) {
		// new * alpha + old * (1 - alpha)
		this.#current = value * this.#alpha + 
			this.#current * (1 - this.#alpha)
	}

	average() {
		return this.#current;
	}
}
