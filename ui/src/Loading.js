import React from 'react';

import "react-loader-spinner/dist/loader/css/react-spinner-loader.css";
import Loader from 'react-loader-spinner'
import "./Loading.css"

const Loading = () => {
    return (
        <div className="spinner-container">
            <div className="spinner"> 
                <Loader type="Bars" color="#00BFFF" height={80} width={80} /> 
            </div>
        </div>
    )
}

export default Loading;