
function colorAnomalyText(countOfAnomalies){
    if (countOfAnomalies < 7){
        return "green"
    }
    if (countOfAnomalies > 6 && countOfAnomalies < 10){
        return "orange"
    }
    if (countOfAnomalies > 9){
        return "red"
    }
}