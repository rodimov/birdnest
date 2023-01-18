import * as React from "react";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";

export default function SortingTableToolbar() {

    return (
        <Toolbar
            sx={{
                pl: { sm: 2 },
                pr: { xs: 1, sm: 1 },
            }}
        >
            <Typography
                sx={{ flex: '1 1 100%', textAlign: 'center' }}
                variant="h6"
                id="tableTitle"
                component="div"
            >
                The pilots who recently violated the NDZ perimeter (last 10 min)
            </Typography>
        </Toolbar>
    );
}
