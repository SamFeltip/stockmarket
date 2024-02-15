
export class StockCalculator extends HTMLElement {

    //stock-value="{{.Cost}}" stock-count="0" total-cost="0"
    static get observedAttributes() { 
        return [
            "stock-value",
            "stock-count",
        ]
    }

    set stockValue(value){
        this.setAttribute('stock-value', value)
    }

    get stockValue(){
        return parseFloat(this.getAttribute('stock-value'))
    }


    set stockCount(value){
        let rounded_stock_count = Math.floor(value / 1000) * 1000
        
        let input_stock_count = this.querySelector('.input-stock')
        input_stock_count.value = rounded_stock_count

        let input_total_cost = this.querySelector('.input-cash')
        input_total_cost.value = rounded_stock_count * this.stockValue

        return this.setAttribute('stock-count', rounded_stock_count)
    }

    get stockCount(){
        return parseInt(this.getAttribute('stock-count'))
    }

    set totalCost(value){

        let rounded_to_value_stock_count = Math.floor(value / this.stockValue);
        let new_stock_count = Math.floor(rounded_to_value_stock_count / 1000) * 1000
        this.stockCount = new_stock_count

    }
    get totalCost(){
        return this.stockCount * this.stockValue
    }

    constructor() {
        super();
    }

    attributeChangedCallback(name, oldValue, newValue) {}

    connectedCallback() {

        this.querySelectorAll('.rem-stock').forEach(rem_stock_button => rem_stock_button.addEventListener('click', (e) => {
            if(this.stockCount >= 1000) { this.stockCount -= 1000 }
            
        }))

        this.querySelectorAll('.add-stock').forEach(add_stock_button => add_stock_button.addEventListener('click', (e) => {
            this.stockCount += 1000
        }))

        let input_stock_count = this.querySelector('.input-stock')
        input_stock_count.value = this.stockCount
        input_stock_count.addEventListener("change", (e) => {
            this.stockCount = e.target.value
        })

        let input_total_cost = this.querySelector('.input-cash')
        input_total_cost.setAttribute("step", 1000 * this.stockValue)
        input_total_cost.addEventListener("change", (e) => {
            this.totalCost = e.target.value
        })

    }
}

customElements.define('stock-calculator', StockCalculator);