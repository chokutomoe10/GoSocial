import { useNavigate, useParams } from "react-router-dom";
import { API_URL } from "./App";

export const ConfirmationPage = () => {
    const { token } = useParams();
    const redirect = useNavigate();

    const handleConfirm = async () => {const response = await fetch(`${API_URL}/users/activate/${token}`, {
            method: "PUT"
        });
        
        console.log("ini token", token);
        console.log("ini response", response);

        if (response.ok) {
            redirect("/")
        } else {
            alert("Failed to confirm token")
        }
    }

    return (
        <div>
            <h1>Confirmation</h1>
            <button onClick={handleConfirm}>Click to confirm</button>
        </div>
    )
}