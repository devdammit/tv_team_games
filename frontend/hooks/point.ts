import {useEffect, useState} from "react";

export enum PointStatus {
    PAUSED = "PAUSED",
    STARTED = "STARTED",
    FINISHED = "FINISHED",
    UPDATE_RESULT = "UPDATE_RESULT",
}

export type Action = {
    number: number
    results: Result[]
    type: PointStatus
    date: string
}

export type Result = {
    name: string,
    id: string,
    timestamp: string
}

export const useWatchPoints = () => {
    const [startedAt, setStartedAt] = useState<Date | null>()
    const [status, setStatus] = useState<PointStatus>()
    const [results, setResults] = useState<Result[]>([])
    const [roundNum, setRoundNum] = useState<number>()

    useEffect(() => {
        const pointSource = new EventSource(`${process.env.NEXT_PUBLIC_HOST_EVENT_SOURCE}?stream=point`)

        pointSource.addEventListener("message", (event) => {
            const action = JSON.parse(event.data) as Action
            setStatus(action.type)

            if (!action.results) {
                action.results = []
            }

            if (action.type === PointStatus.FINISHED) {
                setResults(action.results.sort((a, b) => {
                    return new Date(a.timestamp).getTime() > new Date(b.timestamp).getTime() ? 1 : -1
                }))
            } else if (action.type === PointStatus.STARTED) {
                setResults([])
                setStartedAt(new Date(action.date))
                setRoundNum(action.number)

            } else if (action.type === PointStatus.UPDATE_RESULT) {

                setResults(action.results.sort((a, b) => {
                    return new Date(a.timestamp).getTime() > new Date(b.timestamp).getTime() ? 1 : -1
                }))
            }
        })

        return () => {
            pointSource.close()
        }
    }, [])

    return {status, results, startedAt, roundNum}
}