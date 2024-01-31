import React, { useState } from 'react';
import axios from 'axios';
import { Activity } from '../types/Activity';

const ActivityTable: React.FC = () => {
    const [activities, setActivities] = useState<Activity[]>([
        {
            "Code": "A",
            "Duration": 3,
            "Dependents": ["T"]
        },
        {
            "Code": "C",
            "Duration": 4,
            "Dependents": ["A", "S"]
        },
        {
            "Code": "E",
            "Duration": 5,
            "Dependents": ["C"]
        },
        {
            "Code": "K",
            "Duration": 11,
            "Dependents": ["A"]
        },
        {
            "Code": "O",
            "Duration": 3,
            "Dependents": []
        },
        {
            "Code": "P",
            "Duration": 6,
            "Dependents": ["C", "S"]
        },
        {
            "Code": "S",
            "Duration": 2,
            "Dependents": ["T"]
        },
        {
            "Code": "T",
            "Duration": 7,
            "Dependents": ["O"]
        },
        {
            "Code": "U",
            "Duration": 3,
            "Dependents": ["E", "K", "P"]
        }
    ]);
    const [respond, setResponse] = useState<string>();

    const handleSend = () => {
        axios.post('http://localhost:8080/pert', activities)
            .then(response => setResponse(response.data?.data))
            .catch(error => console.error('There was an error!', error));
    }

    const handleActivityChange = (index: number, key: keyof Activity, value: string | number | Array<string>) => {
        const updatedActivities = [...activities];
        updatedActivities[index] = { ...updatedActivities[index], [key]: value };
        setActivities(updatedActivities);
    };

    const handleFileChange = (e: React.FormEvent<HTMLInputElement>) => {
        const fileReader = new FileReader();
        const file = e.currentTarget.files;
        if (file && file[0]) {
            fileReader.readAsText(file[0], "UTF-8");
            fileReader.onload = e => {
                const jsonText = e.target?.result;
                if (typeof jsonText === 'string') {
                    const activitiesData: Activity[] = JSON.parse(jsonText);
                    setActivities(activitiesData);
                }
            };
            fileReader.onerror = e => {
                console.error("Error reading file:", e);
            };
        }
    };

    return (
        <>
            <table>
                <thead>
                    <tr>
                        <th>Code</th>
                        <th>Duration</th>
                        <th>Dependents</th>
                    </tr>
                </thead>
                <tbody>
                    {activities.map((activity, index) => (
                        <tr key={activity.Code}>
                            <td>{activity.Code}</td>
                            <td>
                                <input
                                    type="number"
                                    value={activity.Duration}
                                    onChange={(e) => handleActivityChange(index, 'Duration', parseInt(e.target.value))}
                                />
                            </td>
                            <td>
                                <input
                                    type="text"
                                    value={activity.Dependents.join(", ")}
                                    onChange={(e) => handleActivityChange(index, 'Dependents', e.target.value.split(", "))}
                                />
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <input type='file' onChange={handleFileChange} />
            <button onClick={handleSend}>Send</button>
            <p>{JSON.stringify(respond)}</p>
        </>
    );
};

export default ActivityTable;