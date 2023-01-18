import * as React from 'react';

import SortingTable from "./SortingTable";
import {createTableData, TableData} from "./tableData";

import "../styles/PilotsList.css"
import useAPI from "../hooks/useAPI";
import Typography from "@mui/material/Typography";

export default function PilotsList() {
    let rows: TableData[] = []
    let nearest = -1

    const { loading, data } = useAPI(
        window.location.href + "api/pilots"
    );

    if (!loading) {
        for (const pilot of data.data) {
            rows.push(
                createTableData(pilot.id, pilot.first_name, pilot.last_name,
                    pilot.phone, pilot.email, pilot.registration_time,
                    pilot.drone.id, pilot.drone.last_seen, pilot.drone.position_x, pilot.drone.position_y,
                    pilot.drone.distance)
            )
        }

        for (const pilot of rows) {
            let distance = Math.round(pilot.distance / 1000)
            if (nearest === -1 || nearest > distance) {
                nearest = distance
            }
        }
    }

    return (
        <div className={"list-container"}>
            {nearest !== -1 &&
                <Typography
                    sx={{ flex: '1 1 100%', textAlign: 'center' }}
                    variant="h6"
                    id="tableTitle"
                    component="div"
                >
                    The closest confirmed distance to the nest is {nearest}m
                </Typography>
            }
            {loading
                ? <div>Loading pilots info...</div>
                : <SortingTable rows={rows}/>
            }
        </div>
    )
}