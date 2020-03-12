import React from 'react';
import './style.css';
import DiscussionList from '../DiscussionList';

export interface Props {}

function Discussions(props: Props) {
    return (
        <div className="discussions w-full flex justify-center">
            <div className="pl-2 pr-2 px-4 py-2 mx-6 w-full all:w-1/2 md:w-1/2 lg:w-1/2 xl:w-1/3 mx-6">
                <img className="discussions_header-image" src="/logo192.png" alt=""/>
                <div className="title">
                    Convos you can watch
                </div>
                {DiscussionList({})}
                <div className="title pt-4">
                    Convos you can join
                </div>
                {DiscussionList({})}
            </div>
        </div>
    );
}

export default Discussions;