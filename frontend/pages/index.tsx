import Head from 'next/head'
import styles from '../styles/Home.module.css'
import {Card, Skeleton, Statistic, Typography} from "antd";
import {Result, useWatchPoints} from "../hooks/point";
import start from "next/dist/server/lib/start-server";

const positions = [...new Array(10)]


export default function Home() {
    const {status, results, startedAt, roundNum} = useWatchPoints()

    console.log(results)

    return (
        <div className={styles.container}>
            <Head>
                <title>Victo UI</title>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <div className={styles.wrapper}>
                {positions.map((_, index) => {
                    const isWaiting =  !results || typeof results[index] === 'undefined'
                    const result = (!!results && results[index] ? results[index] : {}) as Result
                    let diffMs = 0

                    if (result && startedAt) {
                        const time1 = new Date(result.timestamp).getTime()
                        const time2 = startedAt.getTime()

                        diffMs = Math.round(time1 - time2)
                    }


                    return (
                        <div key={index} className={styles.team}>
                            <Card title={`${index + 1} место`} className={styles.teamCard}>
                                <Skeleton loading={isWaiting} avatar active>
                                    <Card.Meta
                                        title={result.name}
                                        description={`${diffMs} мс`}
                                    />
                                </Skeleton>
                            </Card>
                        </div>
                    )
                })}
            </div>

            <div className={styles.stateBar}>
                <div>
                    <Typography.Text>Состояние раунда:</Typography.Text>
                    <Typography.Text keyboard>{status}</Typography.Text>
                </div>
                <div>
                    <Typography.Text>Раунд №:</Typography.Text>
                    <Typography.Text keyboard>{roundNum}</Typography.Text>
                </div>
            </div>
        </div>
    )
}
