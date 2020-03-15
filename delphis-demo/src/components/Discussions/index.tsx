import React, { useState } from 'react';
import classnames from 'classnames';
import './style.scss';
import DiscussionList from '../DiscussionList';

export interface Props {}

function Discussions(props: Props) {
    const [listSelected, setListSelected] = useState('public');

    return (
        <div className="discussions w-full flex justify-center">
            <div className="pl-2 pr-2 px-4 py-2 mx-6 w-full md:w-8/12 lg:w-8/12 xl:w-8/12 mx-6">    
                <div className="discussion-selector flex justify-left mt-6 ml-4 pb-6">
                    <div className={classnames({selector: true, selected: listSelected === 'public', 'mr-8': true})}>
                        <div onClick={() => setListSelected('public')} className="top-public">
                            Top Public Conversations
                        </div>
                    </div>
                    <div onClick={() => setListSelected('private')} className={classnames({selector: true, selected: listSelected === 'private'})}>
                        <div className="top-private">
                            Top Private Conversations
                        </div>
                    </div>
                </div>
                {DiscussionList({})}
            </div>
        </div>
    );
}

export default Discussions;