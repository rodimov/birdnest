import { useEffect, useState } from "react";

const useAPI = (url: string) => {
    
    const [loading, setLoading] = useState<boolean>(true);
    const [data, setData] = useState<any>(null);

    useEffect(() => {
        const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

        const fetchData = async () => {
            while (true) {
                await fetch(url)
                    .then((response) => {
                        return response.json();
                    })
                    .then((json) => {
                        setLoading(false);
                        setData(json);
                    });
                await sleep(2000)
            }
        }

        fetchData()
    }, [url]);

    return { loading, data };
};

export default useAPI;
