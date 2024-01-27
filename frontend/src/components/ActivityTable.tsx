import React, { useState } from 'react';
import axios from 'axios';
import { Activity } from '../types/Activity';

const ActivityTable: React.FC = () => {
    const [activities, setActivities] = useState<Activity[]>([
        {
            "Code": "A",
            "Duration": 8,
            "Dependents": []
        },
        {
            "Code": "B",
            "Duration": 3,
            "Dependents": ["A"]
        }
    ]);
    const [respond, setResponse] = useState<string>();

    const handleSend = () => {
        axios.post('/api/pert', activities)
            .then(response => setResponse(response.data?.data))
            .catch(error => console.error('There was an error!', error));
    }

    const handleActivityChange = (index: number, key: keyof Activity, value: string | number | Array<string>) => {
        const updatedActivities = [...activities];
        updatedActivities[index] = { ...updatedActivities[index], [key]: value };
        setActivities(updatedActivities);
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
            <button onClick={handleSend}>Send</button>
            <p>{JSON.stringify(respond)}</p>
        </>
    );
};

export default ActivityTable;