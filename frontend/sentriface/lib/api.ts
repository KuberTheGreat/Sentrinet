import axios from "axios";

const API = axios.create({
    baseURL: "https://sentrinet.onrender.com"
})

export default API