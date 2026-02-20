export class RollingAverage {
    #sampleSize;
    #events;
    #count;
    #index;

    constructor(sampleSize = 60) {
        this.#sampleSize = sampleSize;
        this.#events = new Float32Array(sampleSize);
        this.#count = 0;
        this.#index = 0;
    }

    add(event) {
        this.#events[this.#index] = event;
        this.#index = (this.#index + 1) % this.#sampleSize
        this.#count++
    }

    average() {
        let sum = 0;
        for (const dur of this.#events) {
            sum += dur
        }
        return sum / (this.#count < this.#sampleSize
            ? this.#count : this.#sampleSize)
    }
}
