export interface TableData {
    id: string;
    firstName: string;
    lastName: string;
    phone: string;
    email: string;
    registrationTime: string;
    droneId: string;
    lastSeen: string;
    positionX: number;
    positionY: number;
    distance: number;
}

export interface HeadCell {
    disablePadding: boolean;
    id: keyof TableData;
    label: string;
    dataType: TableDataType;
}

export enum TableDataType {
    String = 1,
    Number,
    DateString,
}

export const headCells: readonly HeadCell[] = [
    {
        id: 'id',
        dataType: TableDataType.String,
        disablePadding: true,
        label: 'Pilot ID',
    },
    {
        id: 'firstName',
        dataType: TableDataType.String,
        disablePadding: false,
        label: 'First Name',
    },
    {
        id: 'lastName',
        dataType: TableDataType.String,
        disablePadding: false,
        label: 'Last Name',
    },
    {
        id: 'phone',
        dataType: TableDataType.String,
        disablePadding: false,
        label: 'Phone',
    },
    {
        id: 'email',
        dataType: TableDataType.String,
        disablePadding: false,
        label: 'Email',
    },
    {
        id: 'registrationTime',
        dataType: TableDataType.DateString,
        disablePadding: false,
        label: 'Registration time',
    },
    {
        id: 'droneId',
        dataType: TableDataType.String,
        disablePadding: false,
        label: 'Drone ID',
    },
    {
        id: 'lastSeen',
        dataType: TableDataType.DateString,
        disablePadding: false,
        label: 'Last seen',
    },
    {
        id: 'positionX',
        dataType: TableDataType.Number,
        disablePadding: false,
        label: 'position X',
    },
    {
        id: 'positionY',
        dataType: TableDataType.Number,
        disablePadding: false,
        label: 'position Y',
    },
    {
        id: 'distance',
        dataType: TableDataType.Number,
        disablePadding: false,
        label: 'Distance',
    },
];

export function createTableData(
    id: string,
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
    registrationTime: string,
    droneId: string,
    lastSeen: string,
    positionX: number,
    positionY: number,
    distance: number,
): TableData {
    return {
        id,
        firstName,
        lastName,
        phone,
        email,
        registrationTime,
        droneId,
        lastSeen,
        positionX: Math.round(positionX),
        positionY: Math.round(positionY),
        distance,
    };
}