import { useEffect, useState } from "react";

const useAPI = (url: string) => {
    
    const [loading, setLoading] = useState<boolean>(true);
    const [data, setData] = useState<any>(null);

    useEffect(() => {
        setInterval(() => {
            fetch(url)
                .then((response) => {
                    return response.json();
                })
                .then((json) => {
                    setLoading(false);
                    setData(json);
                });
        }, 10000);
    }, [url]);

    return { loading, data };
};

export default useAPI;
